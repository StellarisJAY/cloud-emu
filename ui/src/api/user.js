import api from "./request";

const user = {
    listUser: async function(query) {
        return await api.get("/user", query);
    },
};

export default user;