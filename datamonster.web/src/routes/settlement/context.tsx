import { Settlement } from "@/api/api";
import { ReactNode, createContext, useState } from "react";

interface Props {
  children?: ReactNode;
}

export const SettlementContext = createContext({
  settlement: {} as Settlement | null,
  setSettlement: (_settlement: Settlement) => {},
});

export const SettlementProvider = ({ children }: Props) => {
  const [settlement, setSettlement] = useState<Settlement | null>(null);

  const value = {
    settlement,
    setSettlement,
  };
  return (
    <SettlementContext.Provider value={value}>
      {children}
    </SettlementContext.Provider>
  );
};
