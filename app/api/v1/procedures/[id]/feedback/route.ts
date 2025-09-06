import { NextRequest } from "next/server"

const RAW = process.env.NEXT_PUBLIC_API_URL || "https://ethio-guide-backend.onrender.com"
const BASE = RAW.replace(/\/$/, "")
const HAS_API_SUFFIX = /\/api\/v1$/.test(BASE)

function backendUrl(path: string) {
	return HAS_API_SUFFIX ? `${BASE}${path}` : `${BASE}/api/v1${path}`
}

export async function GET(req: NextRequest, { params }: { params: { id: string } }) {
	const { id } = params
	const search = req.nextUrl.search || ""
	const dest = backendUrl(`/procedures/${encodeURIComponent(id)}/feedback${search}`)
	const headers: Record<string, string> = {
		Accept: req.headers.get("accept") || "application/json",
	}
	const auth = req.headers.get("authorization")
	if (auth) headers["authorization"] = auth
	const res = await fetch(dest, { headers, cache: "no-store" })
	const body = await res.arrayBuffer()
	return new Response(body, {
		status: res.status,
		headers: {
			"content-type": res.headers.get("content-type") || "application/json",
			"cache-control": "no-store",
		},
	})
}

type FeedbackPayload = { Content: string; Type: string; Tags?: string[]; ProcedureID?: string }

export async function POST(req: NextRequest, { params }: { params: { id: string } }) {
	const { id } = params
	const dest = backendUrl(`/procedures/${encodeURIComponent(id)}/feedback`)

	// Accept incoming payload and convert to backend-required keys when needed.
	let payload: FeedbackPayload | Record<string, unknown> | null = null
	try {
		const json: unknown = await req.json()
		const hasBackendKeys = !!(json && typeof json === 'object' && (Object.prototype.hasOwnProperty.call(json, 'Content') || Object.prototype.hasOwnProperty.call(json, 'Type')))
		if (hasBackendKeys) {
			// Already in the expected format; normalize Type
			const norm = (val: unknown) => String(val ?? '').toLowerCase().trim().replace(/\s+/g,'_')
			const toEnum = (v: string) => {
				const s = norm(v)
				if (['inaccuracy','inacuuracy','inacuracy','incorrect','error','issue'].includes(s)) return 'inaccuracy'
				if (['improvement','inmprovement','improvment','enhancement','suggestion'].includes(s)) return 'improvement'
				return 'other'
			}
			payload = { ...(json as Record<string, unknown>), Type: toEnum((json as Record<string, unknown>).Type as string) }
		} else {
			// Build { Content, Type, Tags? } from common lower-case inputs
			const obj = (json && typeof json === 'object') ? (json as Record<string, unknown>) : {}
			const contentStr = String((obj.content ?? obj.body ?? obj.message ?? '') as string).trim()
			const rawType = String((obj.type ?? obj.feedbackType ?? 'inaccuracy') as string)
			const mapToEnum = (v: string) => {
				const s = String(v || '').toLowerCase().trim().replace(/\s+/g,'_')
				if (['inaccuracy','inacuuracy','inacuracy','incorrect','error','issue'].includes(s)) return 'inaccuracy'
				if (['improvement','inmprovement','improvment','enhancement','suggestion'].includes(s)) return 'improvement'
				return 'other'
			}
			const maybeTags = Array.isArray(obj.tags) ? obj.tags : (Array.isArray(obj.Tags as unknown[]) ? (obj.Tags as string[]) : undefined)
			const tags = Array.isArray(maybeTags)
				? (maybeTags as unknown[]).filter((t) => typeof t === 'string' && String(t).trim()).slice(0, 5) as string[]
				: undefined
			payload = { Content: contentStr, Type: mapToEnum(rawType), ...(tags && tags.length ? { Tags: tags } : {}), ProcedureID: id }
		}
	} catch {
		payload = { Content: '', Type: 'GENERAL' }
	}

	const headers: Record<string, string> = {
		"Content-Type": "application/json",
		Accept: req.headers.get("accept") || "application/json",
	}
	const auth = req.headers.get("authorization")
	if (auth) headers["authorization"] = auth

	// Forward to backend
	let res = await fetch(dest, {
		method: 'POST',
		headers,
		body: JSON.stringify(payload),
		cache: 'no-store'
	})

	// If the backend rejects due to Type oneof validation, try alternative enum candidates
	if (res.status === 400 || res.status === 422) {
		let errorText = ''
		try { errorText = await res.clone().text() } catch {}
		const mentionsOneOf = /oneof/i.test(errorText) && /Type/i.test(errorText)
		if (mentionsOneOf) {
			const p = payload as Record<string, unknown>
			const contentValue = String((p.Content ?? (p as any).content ?? '') as string)
			const tagsValue = Array.isArray((p as any).Tags) ? (p as any).Tags as string[] : (Array.isArray((p as any).tags) ? (p as any).tags as string[] : undefined)
			const extracted: string[] = []
			// Try to extract allowed values from messages like: oneof=INACCURACY IMPROVEMENT GENERAL
			const oneofMatch = errorText.match(/oneof[^A-Za-z0-9_\-]*([A-Za-z0-9_\- ,|/]+)/i)
			if (oneofMatch && oneofMatch[1]) {
				oneofMatch[1].split(/[\s,|/]+/).forEach(t => {
					const v = t.trim().replace(/["'\[\]\(\):]+/g, '')
					if (v && v.length <= 40) extracted.push(v)
				})
			}
			// Also parse bracketed lists: [...]
			const bracket = errorText.match(/\[([^\]]+)\]/)
			if (bracket && bracket[1]) {
				bracket[1].split(/[\s,|/]+/).forEach(t => {
					const v = t.trim().replace(/["'\[\]\(\):]+/g, '')
					if (v && !extracted.includes(v) && v.length <= 40) extracted.push(v)
				})
			}
			// Fallback candidate set if extraction fails
			const defaults = ['inaccuracy','improvement','other','general','feature_request']
			const candidates = (extracted.length ? extracted : defaults)
			const tried: string[] = []
			for (const raw of candidates) {
				const lower = raw.toLowerCase()
				const snake = lower.replace(/\s+/g,'_')
				const upper = raw.toUpperCase()
				const upperSnake = upper.replace(/\s+/g,'_')
				const title = raw.charAt(0).toUpperCase() + raw.slice(1).toLowerCase()
				const variants = Array.from(new Set([
					lower, snake, upper, upperSnake, title
				]))
				let success = false
				for (const cand of variants) {
					tried.push(cand)
					const attemptPayload: FeedbackPayload = { Content: contentValue, Type: cand, ProcedureID: id, ...(tagsValue && tagsValue.length ? { Tags: tagsValue } : {}) }
					const attempt = await fetch(dest, {
						method: 'POST',
						headers,
						body: JSON.stringify(attemptPayload),
						cache: 'no-store'
					})
					if (attempt.ok) { res = attempt; success = true; break }
				}
				if (success) break
			}
			if (!res.ok) {
				// Attach tried variants to the error details to aid debugging
				return new Response(JSON.stringify({ error: 'Type oneof mismatch', details: { tried } }), { status: res.status, headers: { 'content-type': 'application/json', 'cache-control': 'no-store' } })
			}
		}
	}


		if (!res.ok) {
			let text = ''
			try { text = await res.text() } catch {}
			console.error('[feedback POST] backend error:', res.status, text)
			let errorMsg = 'Bad Request'
			let details: unknown = undefined
			try {
				const parsed: unknown = text ? JSON.parse(text) : null
				if (parsed && typeof parsed === 'object') {
					errorMsg = ((parsed as Record<string, unknown>).error as string) || ((parsed as Record<string, unknown>).message as string) || errorMsg
					details = parsed
				} else if (typeof parsed === 'string' && parsed) {
					errorMsg = parsed
				}
			} catch {
				if (text) errorMsg = text
			}
			return new Response(JSON.stringify({ error: errorMsg, details }), {
				status: res.status,
				headers: { 'content-type': 'application/json', 'cache-control': 'no-store' },
			})
		}

		const bodyBuf = await res.arrayBuffer()
		return new Response(bodyBuf, {
			status: res.status,
			headers: {
				"content-type": res.headers.get("content-type") || "application/json",
				"cache-control": "no-store",
			},
		})
}

