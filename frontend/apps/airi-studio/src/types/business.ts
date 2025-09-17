export interface Project {
    id: string;
    name: string;
    description: string;
    status: 'active' | 'inactive' | 'archived';
    createdAt: string;
    updatedAt: string;
}

export interface Agent {
    id: string;
    name: string;
    description: string;
    model: string;
    status: 'running' | 'stopped' | 'error';
    createdAt: string;
    updatedAt: string;
}

export interface User {
    id: string;
    name: string;
    email: string;
    avatar?: string;
}