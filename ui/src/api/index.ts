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

    public async saveDownload(download: DownloadInput): Promise<void> {
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

type RawDownload = {
    id: string;
    source: string;
    name: string;
    createdAt: string;
    totalBytes: number;
    completeBytes: number;
    files: string[];
};

export class DownloadInput {
    public id: string;
    public source: string;

    constructor({ source }: { source: string }) {
        this.id = uuid4();
        this.source = source;
    }
}

export class Download {
    public id: string;
    public source: string;
    public name: string;
    public percentDone: number;
    public files: string[];

    constructor({
        id,
        source,
        name,
        completeBytes,
        totalBytes,
        files
    }: RawDownload) {
        this.id = id;
        this.name = name;
        this.source = source;
        this.percentDone = Math.round((completeBytes / totalBytes) * 100);
        this.files = files;
    }
}

export type Profile = { isActive: boolean };
