/* eslint-disable @typescript-eslint/no-use-before-define */
/* eslint-disable @typescript-eslint/explicit-function-return-type */
/* eslint-disable @typescript-eslint/no-explicit-any */
import { v4 as uuid4 } from "uuid";
import Result from "./Result";

type RawApiResponse<P extends unknown, E extends string> =
    | { status: "success"; payload: P; error: undefined }
    | { status: "error"; payload: undefined; error: E };

export default class API {
    constructor(private baseURL: string) {}

    public async saveDownload(download: Download): Promise<void> {
        await this.request("post", "/downloads", download);
    }

    async getDownloads() {
        const resp = await this.request<RawDownload[]>("get", "/downloads");
        if (resp.status === "error") {
            return Result.fail(resp.error);
        }
        return Result.ok(resp.payload.map(p => new Download(p)));
    }

    public async getProfile() {
        const resp = await this.request("get", "/auth/profile");
        if (resp.status === "error") {
            return Result.fail(resp.error);
        }
        return Result.ok(resp.payload as Profile);
    }

    public async login(password: string) {
        const resp = await this.request("post", "/auth/login", { password });
        if (resp.status === "error") {
            return Result.fail(resp.error);
        }
        return Result.ok();
    }

    public async logout(): Promise<void> {
        await this.request("post", "/auth/logout");
    }

    private async request<P, E extends string = string>(
        method: "get" | "post",
        path: string,
        payload?: unknown
    ): Promise<RawApiResponse<P, E>> {
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

type RawDownload = { id: string; source: string };

export class Download {
    public source: string;
    public id: string;

    constructor({ id, source }: { id?: string; source: string }) {
        this.id = id || uuid4();
        this.source = source;
    }
}

export type Profile = { isActive: boolean };