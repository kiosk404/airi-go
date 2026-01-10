import { httpClient } from '@/services/core';

const API_BASE = '/api/foundation/v1/users';

import type { ILoginByPasswordRequestArgs } from "@/api/generated/foundation/user/LoginByPasswordRequest";
import type { ILoginByPasswordResponseArgs } from "@/api/generated/foundation/user/LoginByPasswordResponse";
import type { IUserRegisterRequestArgs } from "@/api/generated/foundation/user/UserRegisterRequest";
import type { IUserRegisterResponseArgs } from "@/api/generated/foundation/user/UserRegisterResponse";
import type { IGetUserInfoResponseArgs } from "@/api/generated/foundation/user/GetUserInfoResponse";
import type { IUserInfoDetailArgs} from "@/api/generated/foundation/domain/user/UserInfoDetail";

export type UserInfo = IUserInfoDetailArgs
export type LoginRequest = ILoginByPasswordRequestArgs
export type LoginResponse = ILoginByPasswordResponseArgs
export type RegisterRequest = IUserRegisterRequestArgs
export type RegisterResponse = IUserRegisterResponseArgs
export type GetUserInfoResponse = IGetUserInfoResponseArgs

// ============ API 方法 ============

/**
 * 密码登录
 * POST /api/foundation/v1/users/login
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
