import api from "./request";

const room = {
    listJoinedRooms: async function(query) {
        return await api.get("/rooms/joined", query);
    },
    listAllRooms: async function(query) {
        return await api.get("/rooms", query);
    }
};

export default room;