import axios from "axios";
import tokenStorage from "./token";
import router from "../router/index"
import configs from "./const";

const v1 = axios.create({
    baseURL: configs.ApiServer,
});

v1.interceptors.request.use(config => {
    const token = tokenStorage.getToken()
    if (token) {
        config.headers.set("Authorization", "Bearer " + token)
    }
    return config
});

v1.interceptors.response.use(r => {
    if (r.headers && r.headers.getAuthorization()) {
        tokenStorage.setToken(r.headers.getAuthorization());
    }
    const resp = r["data"];
    if (resp && resp.code === 200) {
        return resp;
    }
    if (resp && resp.code === 401) return router.push("/login");
    return Promise.reject(resp);
});

function errorHandler(err) {
    if (err.response && err.response.status === 401) return router.push("/login");
    if (err["code"] === 401) {
        return router.push("/login");
    }
    return Promise.reject(err);
}

const api = {
    axios: v1,
    get(path, queryParams) {
        if (queryParams) {
            let params = [];
            for (let key in queryParams) {
                params.push(key + "=" + queryParams[key]);
            }
            path += "?" + params.join("&");
        }
        return this.axios.get(path).catch(err => errorHandler(err));
    },
    post(path, data) {
        return this.axios.post(path, data).catch(err => errorHandler(err));
    },
    put(path, data) {
        return this.axios.put(path, data).catch(err => errorHandler(err));
    },
    delete(path) {
        return this.axios.delete(path).catch(err => errorHandler(err));
    }
}

export default api