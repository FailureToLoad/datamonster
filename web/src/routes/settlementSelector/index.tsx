import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import api, { Settlement } from "@/api/settlement";
import { useCallback, useState } from "react";
import { CreateSettlementDialogue } from "./createSettlementDialogue";
import { Button } from "@/components/ui/button";
import { Play } from "lucide-react";
import { useLoaderData, useNavigate } from "react-router-dom";

function SettlementCard({ settlement }: { settlement: Settlement }) {
  const navigate = useNavigate();
  const navigateOnClick = () => {
    localStorage.setItem("settlement", JSON.stringify(settlement));
    navigate(`..`);
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>{settlement.name}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex flex-row justify-between">
          <div>Lantern Year: {settlement.year}</div>
          <div>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => navigateOnClick()}
            >
              <Play className="h-6 w-6" onClick={() => navigateOnClick()} />
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

interface SettlementListProps {
  update: (s: Settlement) => void;
  settlements: Array<Settlement>;
}

function SettlementList({ update, settlements }: SettlementListProps) {
  const dialogueListItem = (
    <li key={-1}>
      <CreateSettlementDialogue update={update} />
    </li>
  );
  if (settlements === null) {
    return dialogueListItem;
  }

  const cards = settlements.map((settlement) => (
    <li key={settlement.id}>
      <SettlementCard settlement={settlement} />
    </li>
  ));
  return [dialogueListItem, ...cards];
}

export async function settlementListLoader() {
  return await api.getSettlementsForUser();
}

function SettlementSelector() {
  const loadedSettlements = useLoaderData() as Array<Settlement>;
  const [settlements, setSettlements] =
    useState<Array<Settlement>>(loadedSettlements);
  const updateSettlementList = useCallback(
    (s: Settlement) => {
      if (settlements.length === 0) {
        setSettlements([s]);
        return;
      }
      setSettlements([s, ...settlements]);
    },
    [settlements],
  );

  return (
    <div className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden">
      <ul className="w-1/4 space-y-4 ">
        <SettlementList
          update={updateSettlementList}
          settlements={settlements}
        />
      </ul>
    </div>
  );
}

export default SettlementSelector;
