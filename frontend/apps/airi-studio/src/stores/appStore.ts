import { create } from 'zustand';

interface AppState {
    sidebarCollapsed: boolean;
    theme: 'light' | 'dark';
    loading: boolean;
    toggleSidebar: () => void;
    setTheme: (theme: 'light' | 'dark') => void;
    setLoading: (loading: boolean) => void;
}

export const useAppStore = create<AppState>((set) => ({
    sidebarCollapsed: false,
    theme: 'light',
    loading: false,
    toggleSidebar: () => set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),
    setTheme: (theme) => set({ theme }),
    setLoading: (loading) => set({ loading }),
}));
