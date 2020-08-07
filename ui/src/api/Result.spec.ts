import Result from "./Result";

type User = { name: string };

describe("Successful result", () => {
    const result = Result.ok<User>({ name: "Joe" });

    it("is ok", () => {
        expect(result.isOk()).toBe(true);
    });

    it("is not error", () => {
        expect(result.isError()).toBe(false);
    });

    it("allows accessing value", () => {
        const user = result.getValue();
        expect(user).toEqual({ name: "Joe" });
    });

    it("does not allow accessing error", () => {
        expect(() => result.getError()).toThrowErrorMatchingInlineSnapshot(
            `"Unable to get error of successfull result"`
        );
    });
});

describe("Failed result", () => {
    const result = Result.fail({ message: "not_found" });

    it("is not ok", () => {
        expect(result.isOk()).toBe(false);
    });

    it("is error", () => {
        expect(result.isError()).toBe(true);
    });

    it("allows accessing error", () => {
        const error = result.getError();
        expect(error).toEqual({ message: "not_found" });
    });

    it("does not allow accessing value", () => {
        expect(() => result.getValue()).toThrowErrorMatchingInlineSnapshot(
            `"Unable to get value of failed result"`
        );
    });
});

test.skip("typings sensibility", () => {
    // Left here as some kind of a sandbox
    // to play around and have a feel for how typings work
    const fn = (a: number) => {
        if (a === 1) {
            return Result.ok<User>({ name: "Joe" });
        }
        if (a === 2) {
            return Result.fail({ message: "2" as const, two: "hello from 2" });
        }
        if (a === 3) {
            return Result.fail({ message: "3" as const, three: "hi from 3" });
        }
        throw new Error("Unexpected input");
    };

    const okResult = fn(1);
    if (okResult.isOk()) {
        const u = okResult.getValue();
        console.log(u.name);
    }

    const errResult = fn(2);
    if (errResult.isError()) {
        const err = errResult.getError();
        console.log(err.message);
        if (err.message === "2") {
            console.log(err.two);
        }
        if (err.message === "3") {
            console.log(err.three);
        }
    }
});
