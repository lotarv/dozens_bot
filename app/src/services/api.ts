import axios, { AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from "axios";
import { getTelegramInitData } from "./auth";

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    withCredentials: true,
});

const pendingRequests = new Map<string, boolean>();

const getRequestKey = (config: AxiosRequestConfig): string => {
    const {url, method, params} = config
    return `${method?.toUpperCase()}|${url}|${JSON.stringify(params)}`;
}

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
    const initData = getTelegramInitData();
    if (initData) {
        config.headers['X-Telegram-Init-Data'] = initData;
    }

    if (config.method?.toUpperCase() !== "GET") {
        const key = getRequestKey(config)
        if (pendingRequests.has(key)) {
            console.warn(`repeated request rejected: ${key}`)
            return Promise.reject(new axios.Cancel(`Duplicate of request: ${key}`))
        }
        pendingRequests.set(key, true)
    }
    return config
})

api.interceptors.response.use(
    (response:AxiosResponse) => {
        const key = getRequestKey(response.config);
        pendingRequests.delete(key);
        return response;
    },
    (error) => {
        const key = getRequestKey(error.config || {});
        pendingRequests.delete(key)

        if (axios.isCancel(error)) {
            console.warn("Request rejected: ", error.message);
        } else if (error.response?.status == 401) {
            console.error("(401) Authorization error: ", error.message)
        }
        return Promise.reject(error)
    }
)

export default api