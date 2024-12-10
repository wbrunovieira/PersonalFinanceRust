// This file was bundled by Encore v1.45.0
//
// https://encore.dev

// encore.gen/internal/entrypoints/combined/main.ts
import { registerGateways, registerHandlers, run } from "encore.dev/internal/codegen/appinit";

// url/url.ts
import { api, APIError } from "encore.dev/api";
import { SQLDatabase } from "encore.dev/storage/sqldb";
import { randomBytes } from "crypto";
var db = new SQLDatabase("url", { migrations: "./migrations" });
var shorten = api(
  { expose: true, auth: false, method: "POST", path: "/url" },
  async ({ url }) => {
    const id = randomBytes(6).toString("base64url");
    await db.exec`
        INSERT INTO url (id, original_url)
        VALUES (${id}, ${url})
    `;
    return { id, url };
  }
);
var get = api(
  { expose: true, auth: false, method: "GET", path: "/url/:id" },
  async ({ id }) => {
    const row = await db.queryRow`
        SELECT original_url FROM url WHERE id = ${id}
    `;
    if (!row)
      throw APIError.notFound("url not found");
    return { id, url: row.original_url };
  }
);
var list = api(
  { expose: false, method: "GET", path: "/url" },
  async () => {
    const rows = db.query`
        SELECT id, original_url
        FROM url
    `;
    const urls = [];
    for await (const row of rows) {
      urls.push({ id: row.id, url: row.original_url });
    }
    return { urls };
  }
);

// url/encore.service.ts
import { Service } from "encore.dev/service";
var encore_service_default = new Service("url");

// encore.gen/internal/entrypoints/combined/main.ts
var gateways = [];
var handlers = [
  {
    apiRoute: {
      service: "url",
      name: "shorten",
      handler: shorten,
      raw: false,
      streamingRequest: false,
      streamingResponse: false
    },
    endpointOptions: { "auth": false, "expose": true, "isRaw": false, "isStream": false },
    middlewares: encore_service_default.cfg.middlewares || []
  },
  {
    apiRoute: {
      service: "url",
      name: "get",
      handler: get,
      raw: false,
      streamingRequest: false,
      streamingResponse: false
    },
    endpointOptions: { "auth": false, "expose": true, "isRaw": false, "isStream": false },
    middlewares: encore_service_default.cfg.middlewares || []
  },
  {
    apiRoute: {
      service: "url",
      name: "list",
      handler: list,
      raw: false,
      streamingRequest: false,
      streamingResponse: false
    },
    endpointOptions: { "auth": false, "expose": false, "isRaw": false, "isStream": false },
    middlewares: encore_service_default.cfg.middlewares || []
  }
];
registerGateways(gateways);
registerHandlers(handlers);
await run();
//# sourceMappingURL=main.mjs.map
