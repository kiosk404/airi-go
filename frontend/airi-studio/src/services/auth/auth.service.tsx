import * as apiClient from './api-client';
import type {UserInfo, LoginRequest, RegisterRequest} from "./api-client";

export type { UserInfo };

const USER_INFO_KEY = 'user_info';

//======= Service 方法

// 登录
export async function login(account: string, password: string): Promise<UserInfo | null> {
    const req: LoginRequest = {
        account,
        password,
    };
    const resp = await apiClient.loginByPassword(req);

    if (resp.user_info) {
        localStorage.setItem(USER_INFO_KEY, JSON.stringify(resp.user_info));
    }

    return resp.user_info || null;
}


// 注册
export async function register(account: string, password: string, locale?: string): Promise<UserInfo | null> {
    const req: RegisterRequest = {
        account,
        password,
        locale,
    };
    const resp = await apiClient.register(req);

    if (resp.user_info) {
        localStorage.setItem(USER_INFO_KEY, JSON.stringify(resp.user_info));
    }

    return resp.user_info || null;
}

// 退出登录
export async function logout(): Promise<void> {
   try {
       await apiClient.logout();
   } finally {
       localStorage.removeItem(USER_INFO_KEY);
   }
}


export async function getCurrentUser(): Promise<UserInfo | null> {
   try {
       const resp = await apiClient.getUserInfoByToken();
       if (resp.user_info) {
           localStorage.setItem(USER_INFO_KEY, JSON.stringify(resp.user_info));
           return resp.user_info
       }
       return null;
   } catch {
       localStorage.removeItem(USER_INFO_KEY);
       return null;
   }
}

export function getStoredUser(): UserInfo | null {
    const stored = localStorage.getItem(USER_INFO_KEY);
    if (stored) {
        try {
            return JSON.parse(stored);
        } catch {
            return null;
        }
    }
    return null
}

// 检查是否已登录
export function isLoggedIn(): boolean {
    return getStoredUser() !== null;
}