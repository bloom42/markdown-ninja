import countries from "mdninja-js/src/data/countries";

let countriesByCode: Map<string, string> | null = null;

export function countryName(countryCode: string): string {
  if (!countriesByCode) {
    countriesByCode = new Map();
    for (let country of countries) {
      countriesByCode.set(country.code, country.name);
    }
  }

  return countriesByCode.get(countryCode) ?? 'Unknown';
}
