const API_URL = import.meta.env.VITE_API_URL;

function getHeaders(options: RequestInit): HeadersInit {
	const headers: Record<string, string> = {
		"Content-Type": "application/json",
	};

	const telegram = window.Telegram?.WebApp;

	if (telegram?.initData) {
		headers["X-Telegram-Init-Data"] = telegram.initData;
	}

	if (options.headers) {
		Object.assign(headers, options.headers);
	}

	return headers;
}

export async function apiRequest<T>(
	path: string,
	options: RequestInit = {},
): Promise<T> {
	const response = await fetch(`${API_URL}${path}`, {
		...options,
		headers: getHeaders(options),
	});

	if (!response.ok) {
		const text = await response.text();

		throw new Error(
			`API error ${response.status}: ${text}`,
		);
	}

	return response.json();
}