import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import api from "@/api/api";
import { Settlement } from "@/api/settlement";
import { CreateSettlementDialogue } from "./createSettlementDialogue";
import { Button } from "@/components/ui/button";
import { Play } from "lucide-react";
import { Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { Auth0ContextInterface, User, useAuth0 } from "@auth0/auth0-react";
import Spinner from "@/components/spinner";

function SettlementCard({ settlement }: { settlement: Settlement }) {
  const link = "/" + settlement.id;
  return (
    <Card>
      <CardHeader>
        <CardTitle>{settlement.name}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex flex-row justify-between">
          <div>Lantern Year: {settlement.year}</div>
          <div>
            <Link to={link}>
              <Button variant="ghost" size="icon">
                <Play className="h-6 w-6" />
              </Button>
            </Link>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

interface SettlementListProps {
  settlements: Array<Settlement>;
}

function SettlementList({ settlements }: SettlementListProps) {
  const dialogueListItem = (
    <li key={-1}>
      <CreateSettlementDialogue />
    </li>
  );
  if (settlements === null || settlements.length === 0) {
    return dialogueListItem;
  }

  console.log(settlements);

  const cards = settlements.map((settlement) => (
    <li key={settlement.id}>
      <SettlementCard settlement={settlement} />
    </li>
  ));
  return [dialogueListItem, ...cards];
}

async function getSettlementsForUser(
  client: Auth0ContextInterface<User>,
): Promise<Array<Settlement> | null> {
  const token = await client.getAccessTokenWithPopup({
    authorizationParams: {
      audience: import.meta.env.VITE_AUTH0_AUDIENCE,
      scope: "read:settlements",
    },
  });
  const config = {
    headers: { Authorization: `Bearer ${token}` },
  };
  try {
    const response = await api.get<Array<Settlement> | null>(
      `http://localhost:8080/settlement`,
      config,
    );
    if (!response.data) {
      return null;
    }
    return response.data;
  } catch (e) {
    console.log(e);
    return null;
  }
}

type SelectorProps = {
  testId?: string;
};

function SettlementSelector({ testId }: SelectorProps) {
  const client = useAuth0();

  const { isPending, isError, data, error } = useQuery({
    queryKey: ["settlements"],
    queryFn: async () => getSettlementsForUser(client),
  });

  if (isPending) {
    return <Spinner />;
  }

  if (isError) {
    return <span>Error: {error.message}</span>;
  }

  return (
    <div className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden">
      <ul className="w-1/4 space-y-4 " data-testid={testId}>
        <SettlementList settlements={data as Array<Settlement>} />
      </ul>
    </div>
  );
}

export default SettlementSelector;
