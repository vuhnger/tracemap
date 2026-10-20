import { LOCATIONS_PATH, ROUTES_DIR } from "./constants/storage";
import { type Locations, locate } from "./locate";
import { trace } from "./trace";

const target = Bun.argv[2];
if (!target) throw new Error("Usage: bun run trace <host>");

const route = await trace(target);
const locationsFile = Bun.file(LOCATIONS_PATH);
const known: Locations = (await locationsFile.exists())
	? await locationsFile.json()
	: {};
const locations = await locate(
	route.hops.flatMap((hop) => hop.ip ?? []),
	known,
);

await Bun.write(locationsFile, JSON.stringify(locations, null, 2));
await Bun.write(
	`${ROUTES_DIR}/${target}-${Date.now()}.json`,
	JSON.stringify(route, null, 2),
);

console.table(
	route.hops.map((hop) => ({
		...hop,
		city: hop.ip ? locations[hop.ip]?.city : undefined,
	})),
);
