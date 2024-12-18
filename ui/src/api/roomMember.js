import api from "./request";

const roomMember = {
    listRoomMember: async function(roomId) {
        return await api.get("/room-member?roomId="+roomId);
    },
    getUserRoomMember: async function(roomId) {
        return await api.get("/room-member/user?roomId="+roomId);
    }
};

export default roomMember;