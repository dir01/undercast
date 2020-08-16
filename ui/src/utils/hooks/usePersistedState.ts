import { useState, useCallback } from "preact/hooks";

const usePersistedState = <S>(
    storageKey: string,
    initialState: S | undefined = undefined
) => {
    const [state, setState] = useState<S>(() => {
        // Pre-render runs in node and has no access to globals available in browsers
        if (typeof window === "undefined") {
            return;
        }
        const serialized = window.localStorage.getItem(storageKey);
        if (!serialized) {
            return initialState;
        }
        try {
            return JSON.parse(serialized);
        } catch {
            return initialState;
        }
    });

    const setStateAndPersist = useCallback(
        (value: S) => {
            setState(value);
            window.localStorage.setItem(storageKey, JSON.stringify(value));
        },
        [storageKey]
    );

    return [state, setStateAndPersist] as const;
};

export default usePersistedState;
