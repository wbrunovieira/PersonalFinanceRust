// This file was bundled by Encore v1.45.0
//
// https://encore.dev

// encore.gen/internal/entrypoints/combined/main.ts
import { registerGateways, registerHandlers, run } from "encore.dev/internal/codegen/appinit";

// hello/hello.ts
import { api } from "encore.dev/api";
var get = api(
  { expose: true, method: "GET", path: "/hello/:name" },
  async ({ name }) => {
    const msg = `Hello ${name}!`;
    return { message: msg };
  }
);

// hello/encore.service.ts
import { Service } from "encore.dev/service";
var encore_service_default = new Service("hello");

// encore.gen/internal/entrypoints/combined/main.ts
var gateways = [];
var handlers = [
  {
    apiRoute: {
      service: "hello",
      name: "get",
      handler: get,
      raw: false,
      streamingRequest: false,
      streamingResponse: false
    },
    endpointOptions: { "auth": false, "expose": true, "isRaw": false, "isStream": false },
    middlewares: encore_service_default.cfg.middlewares || []
  }
];
registerGateways(gateways);
registerHandlers(handlers);
await run();
//# sourceMappingURL=main.mjs.map
