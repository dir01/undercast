export default class Result<T, E = undefined> {
    private value: T | undefined;
    private error: E | undefined;

    constructor({
        value,
        error
    }: {
        value?: T | undefined;
        error?: E | undefined;
    }) {
        this.value = value;
        this.error = error;
    }

    public static ok<T>(value?: T): Result<T> {
        return new Result({ value });
    }

    public static fail<E>(error: E): Result<undefined, E> {
        return new Result({ error });
    }

    public isOk(): this is Result<T, undefined> {
        return !this.error;
    }

    public isError(): this is Result<undefined, E> {
        return Boolean(this.error);
    }

    public getValue(): T {
        if (!this.isOk()) {
            throw new Error("Unable to get value of failed result");
        }
        return this.value as T;
    }

    public getError(): E {
        if (this.isOk()) {
            throw new Error("Unable to get error of successfull result");
        }
        return this.error as E;
    }
}
