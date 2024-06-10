import api from "@/api/settlement";

export { default as Settlement } from "./settlement";

export async function SettlementLoader({ params }: { params: any }) {
  let id = params?.settlementId as string;
  let settlement = await api.getSettlement(id);
  if (!settlement) {
    return null;
  }
  return settlement;
}
