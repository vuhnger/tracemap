import { $ } from "bun";
import { MAX_HOPS, PROBE_WAIT_SECONDS } from "./constants/traceroute";

export type Hop = { ttl: number; ip: string | null; rttMs: number | null };
export type Route = { target: string; startedAt: string; hops: Hop[] };

const HOP_LINE = /^ *(\d+) +(?:(\S+) +([\d.]+) ms|\*)/gm;

export const parseTraceroute = (output: string): Hop[] =>
	[...output.matchAll(HOP_LINE)].map(([, ttl, ip, rtt]) => ({
		ttl: Number(ttl),
		ip: ip ?? null,
		rttMs: rtt ? Number(rtt) : null,
	}));

export const trace = async (target: string): Promise<Route> => {
	const startedAt = new Date().toISOString();
	const output =
		await $`traceroute -I -n -q 1 -w ${PROBE_WAIT_SECONDS} -m ${MAX_HOPS} ${target}`.text();
	return { target, startedAt, hops: parseTraceroute(output) };
};
