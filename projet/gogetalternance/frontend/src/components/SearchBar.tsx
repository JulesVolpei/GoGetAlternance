import { useState } from "react";
import { Search } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { ToggleGroup, ToggleGroupItem } from "./ui/toggle-group.tsx";
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "./ui/popover.tsx";
import { Checkbox } from "./ui/checkbox.tsx";
import { Label } from "./ui/label.tsx";
import { SlidersHorizontal } from "lucide-react";

const sourceWebsites = [
    "Indeed",
    "Welcome to the Jungle",
    "HelloWork",
    "FranceTravail",
];

export interface SearchData {
    query: string;
    type: string;
    sources: string[];
}

interface SearchBarProps {
    compact?: boolean;
    onSearch?: (data: SearchData) => void;
}

export function SearchBar({ compact = false, onSearch }: SearchBarProps) {
    const [query, setQuery] = useState("");
    const [type, setType] = useState("apprenticeship");
    const [selectedSources, setSelectedSources] = useState<string[]>(sourceWebsites);

    const toggleSource = (source: string) => {
        setSelectedSources((prev) =>
            prev.includes(source)
                ? prev.filter((s) => s !== source)
                : [...prev, source]
        );
    };

    const handleSearchClick = () => {
        if (onSearch) {
            onSearch({ query, type, sources: selectedSources });
        }
    };

    return (
        <div className={`w-full max-w-2xl mx-auto space-y-4 ${compact ? "" : ""}`}>
            <div className="flex justify-center">
                <ToggleGroup
                    type="single"
                    value={type}
                    onValueChange={(v) => v && setType(v)}
                    className="border rounded-lg"
                >
                    <ToggleGroupItem value="apprenticeship" className="px-6 data-[state=on]:bg-primary data-[state=on]:text-primary-foreground">
                        Alternance
                    </ToggleGroupItem>
                    <ToggleGroupItem value="internship" className="px-6 data-[state=on]:bg-primary data-[state=on]:text-primary-foreground">
                        Stage
                    </ToggleGroupItem>
                </ToggleGroup>
            </div>

            {/* Search input + filters */}
            <div className="flex gap-2">
                <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                        placeholder="Chercher un domaine ..."
                        value={query}
                        onChange={(e) => setQuery(e.target.value)}
                        className="pl-10"
                    />
                </div>

                <Popover>
                    <PopoverTrigger asChild>
                        <Button variant="outline" size="icon">
                            <SlidersHorizontal className="h-4 w-4" />
                        </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-56" align="end">
                        <div className="space-y-3">
                            <p className="text-sm font-medium">Sources</p>
                            {sourceWebsites.map((source) => (
                                <div key={source} className="flex items-center gap-2">
                                    <Checkbox
                                        id={source}
                                        checked={selectedSources.includes(source)}
                                        onCheckedChange={() => toggleSource(source)}
                                    />
                                    <Label htmlFor={source} className="text-sm font-normal cursor-pointer">
                                        {source}
                                    </Label>
                                </div>
                            ))}
                        </div>
                    </PopoverContent>
                </Popover>

                <Button onClick={handleSearchClick}>
                    <Search className="h-4 w-4 mr-2" />
                    Chercher
                </Button>
            </div>
        </div>
    );
}
