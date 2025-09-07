import { NextRequest, NextResponse } from 'next/server';

export const dynamic = 'force-dynamic';

function splitText(input: string, maxLen = 180): string[] {
	const parts: string[] = [];
	let remaining = input.trim();
	while (remaining.length > 0) {
		if (remaining.length <= maxLen) {
			parts.push(remaining);
			break;
		}
		const slice = remaining.slice(0, maxLen);
		// Prefer split at sentence end or space
		const lastPunct = Math.max(slice.lastIndexOf('. '), slice.lastIndexOf('።'), slice.lastIndexOf('! '), slice.lastIndexOf('? '));
		const lastSpace = slice.lastIndexOf(' ');
		const cut = lastPunct > 40 ? lastPunct + 1 : (lastSpace > 40 ? lastSpace : maxLen);
		parts.push(remaining.slice(0, cut).trim());
		remaining = remaining.slice(cut).trim();
	}
	return parts.filter(Boolean);
}

export async function POST(req: NextRequest) {
	try {
		const { text, lang } = await req.json();
		if (typeof text !== 'string' || text.trim().length === 0) {
			return NextResponse.json({ error: 'Missing text' }, { status: 400 });
		}
		const tl = (typeof lang === 'string' && lang.trim().length > 0 ? lang : 'am').slice(0, 2);
		const chunks = splitText(text, 180);
		const buffers: Buffer[] = [];
		for (const chunk of chunks) {
			const q = encodeURIComponent(chunk);
			const url = `https://translate.google.com/translate_tts?ie=UTF-8&q=${q}&tl=${encodeURIComponent(tl)}&client=tw-ob`;
			const upstream = await fetch(url, {
				headers: {
					'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124 Safari/537.36',
					'Accept': '*/*',
					'Referer': 'https://translate.google.com/',
				},
			});
			if (!upstream.ok) {
				return NextResponse.json({ error: 'TTS upstream failed' }, { status: upstream.status });
			}
			const arrayBuffer = await upstream.arrayBuffer();
			buffers.push(Buffer.from(arrayBuffer));
		}
		const combined = Buffer.concat(buffers);
		return new NextResponse(combined, {
			status: 200,
			headers: {
				'Content-Type': 'audio/mpeg',
				'Cache-Control': 'no-store',
				'Access-Control-Allow-Origin': '*',
			},
		});
	} catch (err) {
		console.error('TTS error', err);
		return NextResponse.json({ error: 'TTS error' }, { status: 500 });
	}
}

export async function OPTIONS() {
	return new NextResponse(null, {
		status: 204,
		headers: {
			'Access-Control-Allow-Origin': '*',
			'Access-Control-Allow-Methods': 'POST, OPTIONS',
			'Access-Control-Allow-Headers': 'Content-Type',
		},
	});
}
