import { IPMAP_LOCATE_URL, IPMAP_TIMEOUT_MS } from "./constants/ipmap";

export type Location = {
	city: string;
	country: string;
	latitude: number;
	longitude: number;
};
export type Locations = Record<string, Location | null>;

type IpmapResponse = {
	location: {
		cityName: string;
		countryCodeAlpha2: string;
		latitude: number;
		longitude: number;
	} | null;
};

const locateIp = async (ip: string): Promise<Location | null> => {
	const response = await fetch(`${IPMAP_LOCATE_URL}/${ip}/best`, {
		signal: AbortSignal.timeout(IPMAP_TIMEOUT_MS),
	});
	if (response.status === 400) return null;
	if (!response.ok)
		throw new Error(`IPmap answered ${response.status} for ${ip}`);
	const { location } = (await response.json()) as IpmapResponse;
	return (
		location && {
			city: location.cityName,
			country: location.countryCodeAlpha2,
			latitude: location.latitude,
			longitude: location.longitude,
		}
	);
};

export const locate = async (
	ips: string[],
	known: Locations,
): Promise<Locations> => {
	const missing = [...new Set(ips)].filter((ip) => !(ip in known));
	const found = await Promise.all(
		missing.map(async (ip) => [ip, await locateIp(ip)] as const),
	);
	return { ...known, ...Object.fromEntries(found) };
};
