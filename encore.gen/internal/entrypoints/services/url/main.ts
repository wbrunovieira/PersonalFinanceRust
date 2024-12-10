import { registerHandlers, run, type Handler } from "encore.dev/internal/codegen/appinit";
import { shorten as shortenImpl0 } from "../../../../../url/url";
import { get as getImpl1 } from "../../../../../url/url";
import { list as listImpl2 } from "../../../../../url/url";
import * as url_service from "../../../../../url/encore.service";

const handlers: Handler[] = [
    {
        apiRoute: {
            service:           "url",
            name:              "shorten",
            handler:           shortenImpl0,
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
            handler:           getImpl1,
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
            handler:           listImpl2,
            raw:               false,
            streamingRequest:  false,
            streamingResponse: false,
        },
        endpointOptions: {"auth":false,"expose":false,"isRaw":false,"isStream":false},
        middlewares: url_service.default.cfg.middlewares || [],
    },
];

registerHandlers(handlers);
await run();
