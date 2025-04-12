export const BASE_URL = process.env.NEXT_PUBLIC_API_URL || "";
export const SERVER_URL = process.env.SERVER_URL || "http://localhost:8080";

export const get = async <T>(url: string, options?: RequestInit): Promise<T> => {
  const response = await fetch(`${BASE_URL}${url}`, options);
  return response.json();
};

export const post = async <T, D>(url: string, data: D, options?: RequestInit): Promise<T> => {
  const response = await fetch(`${BASE_URL}${url}`, {
    ...options,
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
    body: JSON.stringify(data),
  });
  return response.json();
};

export const put = async <T, D>(url: string, data: D, options?: RequestInit): Promise<T> => {
  const response = await fetch(`${BASE_URL}${url}`, {
    ...options,
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
    body: JSON.stringify(data),
  });
  return response.json();
};

export const del = async <T>(url: string, options?: RequestInit): Promise<T> => {
  const response = await fetch(`${BASE_URL}${url}`, options);
  return response.json();
};
