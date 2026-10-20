import { expect, test } from "bun:test";
import { parseTraceroute } from "./trace";

const fixture = (name: string) =>
	Bun.file(`${import.meta.dir}/../fixtures/${name}`).text();

test("parses Linux traceroute output and skips the header", async () => {
	expect(parseTraceroute(await fixture("linux.txt"))).toEqual([
		{ ttl: 1, ip: "5.101.109.1", rttMs: 5.601 },
		{ ttl: 2, ip: "143.244.192.104", rttMs: 3.784 },
		{ ttl: 3, ip: "143.244.224.110", rttMs: 2.453 },
		{ ttl: 4, ip: "143.244.224.107", rttMs: 1.054 },
		{ ttl: 5, ip: "162.158.84.196", rttMs: 16.246 },
		{ ttl: 6, ip: "162.158.84.78", rttMs: 2.345 },
		{ ttl: 7, ip: null, rttMs: null },
		{ ttl: 8, ip: "1.1.1.1", rttMs: 2.641 },
	]);
});

test("parses macOS traceroute output with silent hops", async () => {
	const hops = parseTraceroute(await fixture("macos.txt"));
	expect(hops.map((hop) => hop.ttl)).toEqual([
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	]);
	expect(hops[0]).toEqual({ ttl: 1, ip: "192.168.1.254", rttMs: 4.059 });
	expect(
		hops.slice(1, 10).every((hop) => hop.ip === null && hop.rttMs === null),
	).toBe(true);
	expect(hops[10]).toEqual({ ttl: 11, ip: "1.1.1.1", rttMs: 39.38 });
});
