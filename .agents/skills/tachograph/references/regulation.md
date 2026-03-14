# Regulation Reference

## Key Chapters

Regulatory text has been migrated to the `tachograph` skill references.

- **[03-data-dictionary.md](regulation/03-data-dictionary.md)**: **CRITICAL**. Contains ASN.1 definitions for all data types.
- **[05-tachograph-cards-file-structure.md](regulation/05-tachograph-cards-file-structure.md)**: DF/EF hierarchy and file identifiers.
- **[11-response-message-content.md](regulation/11-response-message-content.md)**: VU data format (TV/TREP).
- **[12-card-downloading.md](regulation/12-card-downloading.md)**: Card data format (TLV).
- **[16-common-security-mechanisms.md](regulation/16-common-security-mechanisms.md)**: Certificates, signatures, and Part A/B mechanisms.

## Related Skills

- **ASN.1**: For understanding BER/DER encoding and notation. See `.agents/skills/asn1/SKILL.md`.

## Project Scope

- **Phase 1 (Current)**: Driver Card & Vehicle Unit (VU).
- **Deferred**: Workshop, Control, Company cards.

## Goals

1. **Compliance**: Full alignment with EU regulation.
2. **Fidelity**: No data loss (Round-trip).
3. **Usability**: High-fidelity Protobuf model.
