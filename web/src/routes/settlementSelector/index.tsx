import api from "@/api/settlement";

export { default as Selector } from "./selector";

export async function SettlementListLoader() {
  return await api.getSettlementsForUser();
}
