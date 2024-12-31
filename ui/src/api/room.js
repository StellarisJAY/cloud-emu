import api from "./request";

const room = {
    listJoinedRooms: async function(query) {
        return await api.get("/rooms/joined", query);
    },
    listAllRooms: async function(query) {
        return await api.get("/rooms", query);
    },
    getRoomDetail: async function(id) {
        return await api.get("/room/" + id);
    },
    updateRoom: async function(room) {
        return await api.put("/room", room);
    }
};

export default room;