import api from "@/api/settlement";
import Selector from "./selector";

export async function SettlementListLoader() {
  return await api.getSettlementsForUser();
}

export default Selector;
