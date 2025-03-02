import type { Asset, GetProjectItem, Project, User, UserAction } from "~/types";
import axios from "axios";

const baseURL: string = useRuntimeConfig().public.apiBase;

export const getProject = async (item: string): Promise<Project> => {
  return axios
    .get(`${baseURL}/project/${item}`)
    .then((response) => response.data as Project)
    .catch((error) => {
      console.error("Error fetching project:", error);
      throw error;
    });
};

export const getProjectDonateUsers = async (
  pid: string
): Promise<UserAction[]> => {
  return axios
    .get(`${baseURL}/donate-users/${pid}`)
    .then((response) => response.data as UserAction[])
    .catch((error) => {
      console.error("Error fetching project donate users:", error);
      throw error;
    });
};

export const getProjects = async (params?: {
  limit?: number;
  offset?: number;
  identity_number?: string;
}): Promise<Project[]> => {
  return axios
    .get(`${baseURL}/projects`, { params })
    .then((response) => response.data as Project[])
    .catch((error) => {
      console.error("Error fetching projects:", error);
      throw error;
    });
};

export const searchUsers = async (params: {
  identity_number?: string;
  prefix?: string;
}): Promise<User[]> => {
  return axios
    .get(`${baseURL}/users/search`, { params })
    .then((response) => response.data as User[])
    .catch((error) => {
      console.error("Error searching users:", error);
      throw error;
    });
};

export const getUserDonateProjects = async (
  identityNumber: string
): Promise<UserAction[]> => {
  return axios
    .get(`${baseURL}/users-donate/${identityNumber}`)
    .then((response) => response.data as UserAction[])
    .catch((error) => {
      console.error("Error fetching user donate projects:", error);
      throw error;
    });
};

export const getUser = async (
  identityNumber: string
): Promise<User | undefined> => {
  return axios
    .get(`${baseURL}/user/${identityNumber}`)
    .then((response) => response.data as User)
    .catch((error) => {
      console.error("Error fetching user:", error);
      throw error;
    });
};

export const getAssets = async (): Promise<Asset[]> => {
  return axios
    .get(`${baseURL}/assets`)
    .then((response) => response.data as Asset[])
    .catch((error) => {
      console.error("Error fetching assets:", error);
      throw error;
    });
};
