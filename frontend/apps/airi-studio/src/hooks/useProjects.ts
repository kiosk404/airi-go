import { useAppStore } from '../stores/useAppStore';
import { useMemo } from 'react';

export const useProjects = () => {
  const { projects, addProject, updateProject, deleteProject } = useAppStore();

  const runningProjects = useMemo(
    () => projects.filter((project) => project.status === 'running'),
    [projects]
  );

  const draftProjects = useMemo(
    () => projects.filter((project) => project.status === 'draft'),
    [projects]
  );

  const getProjectById = (id: string) => {
    return projects.find((project) => project.id === id);
  };

  const getProjectStats = () => {
    return {
      total: projects.length,
      running: runningProjects.length,
      draft: draftProjects.length,
      stopped: projects.filter((p) => p.status === 'stopped').length,
    };
  };

  return {
    projects,
    runningProjects,
    draftProjects,
    getProjectById,
    getProjectStats,
    addProject,
    updateProject,
    deleteProject,
  };
};






