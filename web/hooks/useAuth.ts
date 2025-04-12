'use client';

import { useState, useEffect, useCallback } from 'react';

interface User {
  id: string;
  email: string;
  role: string;
  name?: string;
}

interface AuthState {
  isLoggedIn: boolean;
  token: string | null;
  user: User | null;
  isLoading: boolean;
  error: string | null;
}

export function useAuth() {
  const [auth, setAuth] = useState<AuthState>({
    isLoggedIn: false,
    token: null,
    user: null,
    isLoading: true,
    error: null
  });

  // 初始化时从localStorage读取认证状态并验证token
  useEffect(() => {
    const fetchAuthStatus = async () => {
      const token = localStorage.getItem('auth_token');
      if (!token) {
        setAuth(prev => ({ ...prev, isLoading: false }));
        return;
      }
      
      try {
        // 验证token是否有效
        const response = await fetch('/api/auth/me', {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        
        if (response.ok) {
          const data = await response.json();
          setAuth({
            isLoggedIn: true,
            token: token,
            user: data.user,
            isLoading: false,
            error: null
          });
        } else {
          // Token无效，清除存储
          localStorage.removeItem('auth_token');
          setAuth({
            isLoggedIn: false,
            token: null,
            user: null,
            isLoading: false,
            error: '会话已过期，请重新登录'
          });
        }
      } catch (error) {
        console.error('验证token失败', error);
        localStorage.removeItem('auth_token');
        setAuth({
          isLoggedIn: false,
          token: null,
          user: null,
          isLoading: false,
          error: '验证登录状态时出错'
        });
      }
    };
    
    fetchAuthStatus();
  }, []);

  // 登录函数
  const login = useCallback((token: string) => {
    localStorage.setItem('auth_token', token);
    setAuth(prev => ({
      ...prev,
      isLoggedIn: true,
      token,
      isLoading: true,
      error: null
    }));
    
    // 登录后立即获取用户信息
    fetch('/api/auth/me', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    .then(res => {
      if (res.ok) return res.json();
      throw new Error('获取用户信息失败');
    })
    .then(data => {
      setAuth(prev => ({
        ...prev,
        user: data.user,
        isLoading: false
      }));
    })
    .catch(err => {
      console.error(err);
      setAuth(prev => ({
        ...prev,
        isLoading: false,
        error: '获取用户信息失败'
      }));
    });
  }, []);

  // 登出函数
  const logout = useCallback(() => {
    // 调用后端登出接口（可选）
    fetch('/api/auth/logout')
      .catch(err => console.error('登出请求失败', err))
      .finally(() => {
        localStorage.removeItem('auth_token');
        setAuth({
          isLoggedIn: false,
          token: null,
          user: null,
          isLoading: false,
          error: null
        });
      });
  }, []);

  // 获取认证头信息，用于API请求
  const getAuthHeader = useCallback(() => {
    if (auth.token) {
      return {
        Authorization: `Bearer ${auth.token}`
      };
    }
    return {};
  }, [auth.token]);

  // 清除错误
  const clearError = useCallback(() => {
    setAuth(prev => ({ ...prev, error: null }));
  }, []);

  return {
    isLoggedIn: auth.isLoggedIn,
    token: auth.token,
    user: auth.user,
    isLoading: auth.isLoading,
    error: auth.error,
    login,
    logout,
    getAuthHeader,
    clearError
  };
} 