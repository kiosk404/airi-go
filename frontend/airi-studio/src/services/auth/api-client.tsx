import { httpClient } from '@/services/core';

const API_BASE = '/api/foundation/v1/users';

/** 用户信息 */
export interface UserInfo {
    user_id?: string;
    name?: string;
    nick_name?: string;
    description?: string;
    avatar_uri?: string;
    avatar_url?: string;
    account?: string;
    locale?: string;
    create_time?: number;
    update_time?: number;
}

/** 登录请求 */
export interface LoginRequest {
    account: string;
    password: string;
}

/** 登录响应 */
export interface LoginResponse {
    user_info?: UserInfo;
    token?: string;
    expire_time?: number;
}

/** 注册请求 */
export interface RegisterRequest {
    account: string;
    password: string;
    locale?: string;
}

/** 注册响应 */
export interface RegisterResponse {
    user_info?: UserInfo;
    token?: string;
    expire_time?: number;
}

/** 获取用户信息响应 */
export interface GetUserInfoResponse {
    user_info?: UserInfo;
}

// ============ API 方法 ============

/**
 * 密码登录
 * POST /api/foundation/v1/user/login
 */
export async function loginByPassword(req: LoginRequest): Promise<LoginResponse> {
    return httpClient.post<LoginResponse>(`${API_BASE}/login`, req);
}

/**
 * 用户注册
 * POST /api/foundation/v1/users/register
 */
export async function register(req: RegisterRequest): Promise<RegisterResponse> {
    return httpClient.post<RegisterResponse>(`${API_BASE}/register`, req);
}

/**
 * 退出登录
 * POST /api/foundation/v1/users/logout
 */
export async function logout(): Promise<void> {
    return httpClient.post<void>(`${API_BASE}/logout`);
}

/**
 * 获取当前用户信息（基于 session）
 * GET /api/foundation/v1/users/session
 */
export async function getUserInfoByToken(): Promise<GetUserInfoResponse> {
    return httpClient.get<GetUserInfoResponse>(`${API_BASE}/session`);
}
