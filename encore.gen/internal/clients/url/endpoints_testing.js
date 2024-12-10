import { apiCall, streamIn, streamOut, streamInOut } from "encore.dev/internal/codegen/api";
import { registerTestHandler } from "encore.dev/internal/codegen/appinit";

import * as url_service from "../../../../url/encore.service";

export async function shorten(params) {
    const handler = (await import("../../../../url/url")).shorten;
    registerTestHandler({
        apiRoute: { service: "url", name: "shorten", raw: false, handler, streamingRequest: false, streamingResponse: false },
        middlewares: url_service.default.cfg.middlewares || [],
        endpointOptions: {"auth":false,"expose":true,"isRaw":false,"isStream":false},
    });

    return apiCall("url", "shorten", params);
}

export async function get(params) {
    const handler = (await import("../../../../url/url")).get;
    registerTestHandler({
        apiRoute: { service: "url", name: "get", raw: false, handler, streamingRequest: false, streamingResponse: false },
        middlewares: url_service.default.cfg.middlewares || [],
        endpointOptions: {"auth":false,"expose":true,"isRaw":false,"isStream":false},
    });

    return apiCall("url", "get", params);
}

export async function list(params) {
    const handler = (await import("../../../../url/url")).list;
    registerTestHandler({
        apiRoute: { service: "url", name: "list", raw: false, handler, streamingRequest: false, streamingResponse: false },
        middlewares: url_service.default.cfg.middlewares || [],
        endpointOptions: {"auth":false,"expose":false,"isRaw":false,"isStream":false},
    });

    return apiCall("url", "list", params);
}

