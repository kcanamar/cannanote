import 'package:uuid/uuid.dart';

/// Product type categories
enum ProductType { flower, concentrate, edible, tincture, topical, other }

/// Core domain entity for cannabis products
class Product {
  final String id;
  final String name;
  final String? strain;
  final ProductType type;
  final double? thcPercent;
  final double? cbdPercent;
  final String? notes;
  final DateTime createdAt;
  DateTime updatedAt;

  Product({
    String? id,
    required this.name,
    this.strain,
    required this.type,
    this.thcPercent,
    this.cbdPercent,
    this.notes,
    DateTime? createdAt,
    DateTime? updatedAt,
  })  : id = id ?? const Uuid().v4(),
        createdAt = createdAt ?? DateTime.now(),
        updatedAt = updatedAt ?? DateTime.now();

  /// Domain validation
  void validate() {
    if (name.isEmpty) {
      throw ProductException('Product name is required');
    }
    if (thcPercent != null && (thcPercent! < 0 || thcPercent! > 100)) {
      throw ProductException('THC percent must be between 0 and 100');
    }
    if (cbdPercent != null && (cbdPercent! < 0 || cbdPercent! > 100)) {
      throw ProductException('CBD percent must be between 0 and 100');
    }
  }

  Product copyWith({
    String? id,
    String? name,
    String? strain,
    ProductType? type,
    double? thcPercent,
    double? cbdPercent,
    String? notes,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Product(
      id: id ?? this.id,
      name: name ?? this.name,
      strain: strain ?? this.strain,
      type: type ?? this.type,
      thcPercent: thcPercent ?? this.thcPercent,
      cbdPercent: cbdPercent ?? this.cbdPercent,
      notes: notes ?? this.notes,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}

/// Domain exception for product validation errors
class ProductException implements Exception {
  final String message;
  ProductException(this.message);

  @override
  String toString() => message;
}
