import { afterEach, expect, mock, spyOn, test } from "bun:test";
import { locate } from "./locate";

afterEach(() => mock.restore());

const answers: Record<string, Response> = {
	"192.168.1.254": Response.json({ code: "BadRequest" }, { status: 400 }),
	"84.208.0.1": Response.json({ location: null }),
	"62.115.137.62": Response.json({
		location: {
			cityName: "Los Angeles",
			countryCodeAlpha2: "US",
			latitude: 34.05,
			longitude: -118.24,
			iataCode: "LAX",
		},
	}),
};

test("looks up each unknown IP once and keeps known ones", async () => {
	const fetchSpy = spyOn(globalThis, "fetch").mockImplementation((async (
		url: string,
	) => {
		const ip = url.split("/").at(-2) ?? "";
		return answers[ip] ?? Response.json({}, { status: 503 });
	}) as typeof fetch);
	const oslo = { city: "Oslo", country: "NO", latitude: 59.9, longitude: 10.7 };

	const locations = await locate(
		[
			"1.1.1.1",
			"192.168.1.254",
			"84.208.0.1",
			"62.115.137.62",
			"62.115.137.62",
		],
		{ "1.1.1.1": oslo },
	);

	expect(fetchSpy).toHaveBeenCalledTimes(3);
	expect(locations).toEqual({
		"1.1.1.1": oslo,
		"192.168.1.254": null,
		"84.208.0.1": null,
		"62.115.137.62": {
			city: "Los Angeles",
			country: "US",
			latitude: 34.05,
			longitude: -118.24,
		},
	});
});

test("fails instead of caching a miss when IPmap is down", async () => {
	spyOn(globalThis, "fetch").mockImplementation((async () =>
		Response.json({}, { status: 503 })) as unknown as typeof fetch);

	await expect(locate(["62.115.137.62"], {})).rejects.toThrow("503");
});
