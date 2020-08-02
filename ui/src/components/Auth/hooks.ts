import API, { Profile } from "../../API";
import { useState, useCallback, useEffect } from "preact/hooks";
import { createContainer } from "../../unstated-next-preact";

const useAuth = (
    api: API | undefined
): {
    profile: Profile | null;
    isLoggedIn: boolean;
    isLoading: boolean;
    login: (token: string) => Promise<void>;
    logout: () => Promise<void>;
} => {
    const [isLoggedIn, setLoggedIn] = useState(false);
    const [isLoading, setLoading] = useState(true);
    const [profile, setProfile] = useState<Profile | null>(null);

    useEffect(() => {
        (async (): Promise<void> => {
            if (!api) {
                return;
            }
            const profile = await api.getProfile();
            setProfile(profile);
            setLoading(false);
            setLoggedIn(Boolean(profile));
        })();
    }, [api]);

    const login = useCallback(
        async (password: string) => {
            if (!api) {
                return;
            }
            await api.login(password);
            setLoggedIn(true);
        },
        [api]
    );

    const logout = useCallback(async () => {
        console.log("HELLO");
        if (!api) {
            console.log("NO API");
            return;
        }
        await api.logout();
        setLoggedIn(false);
    }, [api]);

    return { profile, isLoggedIn, isLoading, login, logout };
};

export const AuthContainer = createContainer(useAuth);
