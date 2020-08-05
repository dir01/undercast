/* eslint-disable @typescript-eslint/no-explicit-any */
import { v4 as uuid4 } from "uuid";

export default class API {
    constructor(private baseURL: string) {}

    public async saveDownload(download: Download): Promise<void> {
        await this.request("post", "/downloads", download);
    }

    public async getProfile(): Promise<Profile> {
        const response = await this.request("get", "/auth/profile");
        return response.payload;
    }

    public async login(password: string): Promise<void> {
        await this.request("post", "/auth/login", { password });
    }

    public async logout(): Promise<void> {
        await this.request("post", "/auth/logout");
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

export type Profile = { isActive: boolean };
