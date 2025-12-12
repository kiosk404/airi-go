import { createContext, useContext, useState, useEffect, useCallback, useRef, type ReactNode } from 'react';
import type { UserInfo } from '@/services/auth';
import * as authService from '@/services/auth/auth.service';

interface AuthContextType {
    user: UserInfo | null;
    isLoading: boolean;
    isAuthenticated: boolean;
    login: (account: string, password: string) => Promise<void>;
    register: (account: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
    refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<UserInfo | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const mountedRef = useRef(false);

    // 初始化时检查登录状态
    useEffect(() => {
        if (mountedRef.current) return;
        mountedRef.current = true;

        const initAuth = async () => {
            try {
                // 先从本地存储获取
                const storedUser = authService.getStoredUser();
                if (storedUser) {
                    setUser(storedUser);
                    // 后台验证 session 是否有效
                    const currentUser = await authService.getCurrentUser();
                    if (currentUser) {
                        setUser(currentUser);
                    } else {
                        setUser(null);
                    }
                }
            } catch {
                setUser(null);
            } finally {
                setIsLoading(false);
            }
        };

        initAuth();
    }, []);

    const login = useCallback(async (account: string, password: string) => {
        const userInfo = await authService.login(account, password);
        setUser(userInfo);
    }, []);

    const register = useCallback(async (account: string, password: string) => {
        const userInfo = await authService.register(account, password);
        setUser(userInfo);
    }, []);

    const logout = useCallback(async () => {
        await authService.logout();
        setUser(null);
    }, []);

    const refreshUser = useCallback(async () => {
        const userInfo = await authService.getCurrentUser();
        setUser(userInfo);
    }, []);

    return (
        <AuthContext.Provider
            value={{
                user,
                isLoading,
                isAuthenticated: !!user,
                login,
                register,
                logout,
                refreshUser,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}
