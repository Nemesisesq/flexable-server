var payload = {Payload: {}}
var emp_payload = {EmployeePayload: {}}

payload.Payload.Type = {
    UNKNOWN: 0,
    SHIFT_DETAILS: 1,
    ACCEPT_SHIFT_SUBSTITUE: 2,
    DENY_SHIFT_SUBSTITUTE: 3,
    FIND_SHIFT_SUBSTITUTE: 4,
    GET_AVAILABLE_EMPLOYEES: 5,
    GET_JOBS: 6,
    EMPLOYEE_LIST: 7,
    SELECT_VOLUNTEER: 8,
    OPEN_SHIFTS: 9,
    GET_OPEN_SHIFTS: 10,
    PICK_UP_SHIFT: 11,
    CALL_OFF_SHIFT: 12,
    GET_MY_SCHEDULES: 13,
    USER_PROFILE_UPDATED: 14,
    SET_PROFILE: 15,
    UPDATE_NOTIFICATIONS: 16,
    CLOSEOUT_SHIFT: 17,
    GET_COMPANY_LIST: 18,
};
module.exports = {
    payload: payload,
    employee: emp_payload
}