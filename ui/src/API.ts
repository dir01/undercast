/* eslint-disable @typescript-eslint/no-explicit-any */
import { v4 as uuid4 } from "uuid";

export default class API {
    constructor(private baseURL: string) {}

    public async saveDownload(download: Download): Promise<void> {
        await this.request("post", "/api/downloads", download);
    }

    private async request(
        method: "get" | "post",
        path: string,
        payload?: unknown
    ): Promise<Record<string, any>> {
        const req: any = {
            method,
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json"
            }
        };
        if (payload) {
            req.body = JSON.stringify(payload);
        }
        const response = await fetch(`${this.baseURL}${path}`, req);
        return await response.json();
    }
}

export class Download {
    public source: string;
    public id: string;

    constructor({ id, source }: { id?: string; source: string }) {
        this.id = id || uuid4();
        this.source = source;
    }
}
