import api from "@/api/settlement";
import Selector from "./selector";

export async function SettlementListLoader() {
  let settlements = await api.getSettlementsForUser();
  if (!settlements) {
    return null;
  }
  return settlements;
}

export default Selector;
