import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";

import api, { Settlement } from "@/api/api";
import { useCallback, useEffect, useState } from "react";
import Spinner from "@/components/spinner";
import { CreateSettlementDialogue } from "./createSettlementDialogue";

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
  if (settlements.length < 1) {
    return dialogueListItem;
  }

  const cards = settlements.map((settlement) => (
    <li key={settlement.id}>
      <Card>
        <CardHeader>
          <CardTitle>{settlement.name}</CardTitle>
          <CardDescription>Lantern Year: {settlement.year}</CardDescription>
        </CardHeader>
      </Card>
    </li>
  ));
  return [dialogueListItem, ...cards];
}

function SettlementSelector() {
  const [settlements, setSettlements] = useState<Array<Settlement>>(
    Array<Settlement>(),
  );
  const [isLoading, setIsLoading] = useState(true);
  const updateSettlementList = useCallback(
    (s: Settlement) => {
      setSettlements([s, ...settlements]);
    },
    [settlements],
  );

  useEffect(() => {
    api.getSettlementsForUser().then((val) => {
      setSettlements(val);
      setIsLoading(false);
    });
  }, []);

  return (
    <div className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden">
      {isLoading ? (
        <Spinner />
      ) : (
        <ul className="w-1/4 space-y-4 ">
          <SettlementList
            update={updateSettlementList}
            settlements={settlements}
          />
        </ul>
      )}
    </div>
  );
}

export default SettlementSelector;
