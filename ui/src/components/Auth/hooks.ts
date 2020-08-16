import API, { Profile } from "../../api";
import { useState, useCallback, useEffect } from "preact/hooks";
import { createContainer } from "../../unstated-next-preact";
import usePersistedState from "../../utils/hooks/usePersistedState";

const useAuth = (
    api: API | undefined
): {
    profile: Profile | null;
    isLoggedIn: boolean;
    isLoading: boolean;
    login: (token: string) => Promise<void>;
    logout: () => Promise<void>;
    loginError: string;
} => {
    const [isLoggedIn, setLoggedIn] = usePersistedState("isLoggedIn", false);
    const [isLoading, setLoading] = useState(!isLoggedIn);
    const [profile, setProfile] = usePersistedState<Profile>(
        "profile" as const
    );
    const [loginError, setLoginError] = useState("");

    useEffect(() => {
        (async (): Promise<void> => {
            if (!api) {
                return;
            }
            const profileResult = await api.getProfile();
            setLoading(false);
            if (profileResult.isOk()) {
                const profile = profileResult.getValue();
                setProfile(profile);
                setLoggedIn(true);
            } else {
                setLoggedIn(false);
            }
        })();
    }, [api]);

    const login = useCallback(
        async (password: string) => {
            if (!api) {
                return;
            }
            const result = await api.login(password);
            if (result.isOk()) {
                setLoggedIn(true);
            } else {
                setLoginError(result.getError());
            }
        },
        [api]
    );

    const logout = useCallback(async () => {
        if (!api) {
            return;
        }
        await api.logout();
        setLoggedIn(false);
    }, [api]);

    return { profile, isLoggedIn, isLoading, login, logout, loginError };
};

export const AuthContainer = createContainer(useAuth);