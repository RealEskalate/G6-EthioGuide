import { NextRequest, NextResponse } from 'next/server';
import textToSpeech from '@google-cloud/text-to-speech';

export const dynamic = 'force-dynamic';

const client = new textToSpeech.TextToSpeechClient();

export async function POST(req: NextRequest) {
	try {
		const { text, lang } = await req.json();
		if (typeof text !== 'string' || text.trim().length === 0) {
			return NextResponse.json({ error: 'Missing text' }, { status: 400 });
		}
		const languageCode = (typeof lang === 'string' && lang.toLowerCase().startsWith('am')) ? 'am-ET' : 'am-ET';

		const [response] = await client.synthesizeSpeech({
			input: { text },
			voice: { languageCode, ssmlGender: 'FEMALE' as const },
			audioConfig: { audioEncoding: 'MP3' as const },
		});

		const audio = response.audioContent as string | Uint8Array | null | undefined;
		if (!audio) return NextResponse.json({ error: 'No audio' }, { status: 500 });

		// google SDK returns base64 string for audioContent
		const nodeBuffer = Buffer.isBuffer(audio)
			? audio
			: typeof audio === 'string'
				? Buffer.from(audio, 'base64')
				: Buffer.from(audio);
		const body = new Uint8Array(nodeBuffer);

		return new NextResponse(body, {
			status: 200,
			headers: {
				'Content-Type': 'audio/mpeg',
				'Cache-Control': 'no-store',
			},
		});
	} catch (e) {
		console.error('TTS Error:', e);
		return NextResponse.json({ error: 'TTS failed' }, { status: 500 });
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
