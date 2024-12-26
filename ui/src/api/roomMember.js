import api from "./request";

const roomMember = {
    listRoomMember: async function(roomId) {
        return await api.get("/room-member?roomId="+roomId);
    },
    getUserRoomMember: async function(roomId) {
        return await api.get("/room-member/user?roomId="+roomId);
    },
    inviteRoomMember: async function(roomId, userId) {
        return await api.post("/room-member/invite", { "roomId": roomId, "userId": userId });
    }
};

export default roomMember;