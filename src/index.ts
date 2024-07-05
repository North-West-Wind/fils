import "dotenv/config";
import * as crypto from "crypto";
import express from "express";
import * as fs from "fs";
import * as path from "path";
import { TidyURL } from "tidy-url";

const PASSPHRASE = process.env.PASSPHRASE;
if (!PASSPHRASE)
	console.warn("No passphrase set. Everything can use the link shortener!");
const STORAGE_PATH = process.env.STORAGE_PATH || "runtime/links.json";
if (!fs.existsSync(STORAGE_PATH)) {
	fs.mkdirSync(path.dirname(STORAGE_PATH), { recursive: true });
	fs.writeFileSync(STORAGE_PATH, "{}");
}
const ID_LENGTH = process.env.ID_LENGTH && !isNaN(parseInt(process.env.ID_LENGTH)) ? parseInt(process.env.ID_LENGTH) : 8;
const HASH_ALGORITHM = process.env.HASH_ALGORITHM || "md5";
const WRITE_INTERVAL = process.env.WRITE_INTERVAL && !isNaN(parseInt(process.env.WRITE_INTERVAL)) ? parseInt(process.env.WRITE_INTERVAL) : 60000;

const mapping = JSON.parse(fs.readFileSync(STORAGE_PATH, { encoding: "utf8" }));
TidyURL.config.set("silent", true);

const app = express();
app.use(express.static("public"));
app.use(express.json());

app.get("/", (_req, res) => {
	res.send("index.html");
});

app.get("/s/:id", (req, res) => {
	if (mapping[req.params.id]) res.redirect(301, mapping[req.params.id]);
	else res.sendStatus(404);
});

app.post("/api/new", (req, res) => {
	if (PASSPHRASE && req.headers.authorization !== `Bearer ${PASSPHRASE}`) return res.status(403).send("Wrong passphrase");
	if (!req.body.link) return res.status(400).send("Missing link");
	if (!TidyURL.validate(req.body.link)) return res.status(400).send("Invalid link");
	const clean = TidyURL.clean(req.body.link).url;
	if (!/^https?:$/.test(new URL(clean).protocol)) return res.status(400).send("Link must be either HTTP or HTTPS");
	let input = clean, hash: string, ok = false;
	do {
		hash = crypto.createHash(HASH_ALGORITHM).update(input).digest("hex").slice(0, ID_LENGTH);
		if (!mapping[hash]) {
			mapping[hash] = clean;
			ok = true;
		} else if (mapping[hash] && mapping[hash] == clean) ok = true;
		else input = hash;
	} while (!ok);
	res.send(hash);
});

const server = app.listen(process.env.PORT || 3000, async () => {
	console.log('FILS listening at port %s', (<any>server.address()).port);
});

setInterval(() => {
	fs.writeFileSync(STORAGE_PATH, JSON.stringify(mapping));
}, WRITE_INTERVAL);