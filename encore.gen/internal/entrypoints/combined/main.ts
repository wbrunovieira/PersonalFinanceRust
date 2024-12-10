import { registerGateways, registerHandlers, run, type Handler } from "encore.dev/internal/codegen/appinit";

import { shorten as url_shortenImpl0 } from "../../../../url/url";
import { get as url_getImpl1 } from "../../../../url/url";
import { list as url_listImpl2 } from "../../../../url/url";
import * as url_service from "../../../../url/encore.service";

const gateways: any[] = [
];

const handlers: Handler[] = [
    {
        apiRoute: {
            service:           "url",
            name:              "shorten",
            handler:           url_shortenImpl0,
            raw:               false,
            streamingRequest:  false,
            streamingResponse: false,
        },
        endpointOptions: {"auth":false,"expose":true,"isRaw":false,"isStream":false},
        middlewares: url_service.default.cfg.middlewares || [],
    },
    {
        apiRoute: {
            service:           "url",
            name:              "get",
            handler:           url_getImpl1,
            raw:               false,
            streamingRequest:  false,
            streamingResponse: false,
        },
        endpointOptions: {"auth":false,"expose":true,"isRaw":false,"isStream":false},
        middlewares: url_service.default.cfg.middlewares || [],
    },
    {
        apiRoute: {
            service:           "url",
            name:              "list",
            handler:           url_listImpl2,
            raw:               false,
            streamingRequest:  false,
            streamingResponse: false,
        },
        endpointOptions: {"auth":false,"expose":false,"isRaw":false,"isStream":false},
        middlewares: url_service.default.cfg.middlewares || [],
    },
];

registerGateways(gateways);
registerHandlers(handlers);

await run();
