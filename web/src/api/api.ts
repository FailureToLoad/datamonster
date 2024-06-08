import axios from "axios";
const ax = axios.create({
  headers: {
    "Accept-Language": "en-US,en;q=0.5",
  },
  validateStatus: (status) => {
    return status >= 200 && status < 500;
  },
});
ax.defaults.withCredentials = true;

export default ax;
