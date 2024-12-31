import api from "./request";

const user = {
    listUser: async function(query) {
        return await api.get("/user", query);
    },
    getLoginUserDetail: async function() {
        return await api.get("/login-user");
    },
    activateAccount: async function(code) {
        return await api.post("/account/activate", {code: code});
    }
};

export default user;