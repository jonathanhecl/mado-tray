#import <Cocoa/Cocoa.h>
#import <dispatch/dispatch.h>

extern void madoTrayShow(void);
extern void madoTrayExit(void);

@interface MadoTrayMenuTarget : NSObject
- (void)showWindow:(id)sender;
- (void)exitApp:(id)sender;
@end

@implementation MadoTrayMenuTarget
- (void)showWindow:(id)sender {
  [NSApp activateIgnoringOtherApps:YES];
  madoTrayShow();
}

- (void)exitApp:(id)sender {
  madoTrayExit();
}
@end

static NSStatusItem *madoTrayStatusItem;
static MadoTrayMenuTarget *madoTrayTarget;
static NSMenuItem *madoTrayShowItem;
static NSMenuItem *madoTrayExitItem;
static NSString *madoTrayLocale = @"en";

static NSString *MadoTrayShowLabel(void) {
  if ([madoTrayLocale isEqualToString:@"es"]) {
    return @"Mostrar";
  }
  return @"Show";
}

static NSString *MadoTrayExitLabel(void) {
  if ([madoTrayLocale isEqualToString:@"es"]) {
    return @"Salir de Mado-Tray";
  }
  return @"Exit Mado-Tray";
}

static void MadoTrayStrokeLine(NSPoint from, NSPoint to, CGFloat width) {
  NSBezierPath *line = [NSBezierPath bezierPath];
  [line moveToPoint:from];
  [line lineToPoint:to];
  [line setLineWidth:width];
  [line stroke];
}

static NSImage *MadoTrayIcon(void) {
  const CGFloat size = 18.0;
  NSImage *image = [[NSImage alloc] initWithSize:NSMakeSize(size, size)];

  [image lockFocus];
  [[NSColor blackColor] setStroke];
  [[NSColor blackColor] setFill];

  // Marco exterior de ventana japonesa (mado / shoji).
  NSBezierPath *frame = [NSBezierPath bezierPathWithRect:NSMakeRect(1.5, 2.0, 15.0, 14.0)];
  [frame setLineWidth:1.3];
  [frame stroke];

  // Kamoi: travesaño superior más grueso.
  NSBezierPath *topRail = [NSBezierPath bezierPathWithRect:NSMakeRect(2.0, 13.2, 14.0, 2.2)];
  [topRail fill];

  // Shikii: umbral inferior.
  NSBezierPath *bottomRail = [NSBezierPath bezierPathWithRect:NSMakeRect(2.0, 2.0, 14.0, 1.4)];
  [bottomRail fill];

  // Montantes laterales del enrejado.
  MadoTrayStrokeLine(NSMakePoint(6.2, 3.6), NSMakePoint(6.2, 12.8), 1.0);
  MadoTrayStrokeLine(NSMakePoint(11.8, 3.6), NSMakePoint(11.8, 12.8), 1.0);

  // Travesaños horizontales del shoji (tres filas de paneles).
  MadoTrayStrokeLine(NSMakePoint(2.4, 10.2), NSMakePoint(15.6, 10.2), 0.9);
  MadoTrayStrokeLine(NSMakePoint(2.4, 7.0), NSMakePoint(15.6, 7.0), 0.9);

  [image unlockFocus];
  image.template = YES;
  return image;
}

static void MadoTrayEnsureTarget(void) {
  if (madoTrayTarget == nil) {
    madoTrayTarget = [[MadoTrayMenuTarget alloc] init];
  }
}

static void MadoTrayUpdateMenuTexts(void) {
  if (madoTrayShowItem == nil || madoTrayExitItem == nil) {
    return;
  }

  madoTrayShowItem.title = MadoTrayShowLabel();
  madoTrayExitItem.title = MadoTrayExitLabel();
}

void MadoTrayCreate(void) {
  dispatch_async(dispatch_get_main_queue(), ^{
    MadoTrayEnsureTarget();
  });
}

void MadoTrayShow(void) {
  dispatch_async(dispatch_get_main_queue(), ^{
    MadoTrayEnsureTarget();

    if (madoTrayStatusItem != nil) {
      return;
    }

    madoTrayStatusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSSquareStatusItemLength];
    [madoTrayStatusItem retain];
    madoTrayStatusItem.autosaveName = @"com.jonathanhecl.mado-tray.status-item";
    if (@available(macOS 10.12, *)) {
      madoTrayStatusItem.behavior = 0;
    }
    madoTrayStatusItem.button.toolTip = @"Mado-Tray";

    NSImage *image = MadoTrayIcon();
    if (image != nil) {
      madoTrayStatusItem.button.image = image;
    } else {
      madoTrayStatusItem.button.title = @"Mado";
    }
    madoTrayStatusItem.button.enabled = YES;
    madoTrayStatusItem.button.hidden = NO;

    NSMenu *menu = [[NSMenu alloc] initWithTitle:@"Mado-Tray"];
    madoTrayShowItem = [[NSMenuItem alloc] initWithTitle:MadoTrayShowLabel()
                                                   action:@selector(showWindow:)
                                            keyEquivalent:@""];
    madoTrayShowItem.target = madoTrayTarget;
    [menu addItem:madoTrayShowItem];

    [menu addItem:[NSMenuItem separatorItem]];

    madoTrayExitItem = [[NSMenuItem alloc] initWithTitle:MadoTrayExitLabel()
                                                   action:@selector(exitApp:)
                                            keyEquivalent:@""];
    madoTrayExitItem.target = madoTrayTarget;
    [menu addItem:madoTrayExitItem];

    madoTrayStatusItem.menu = menu;
    madoTrayStatusItem.length = NSSquareStatusItemLength;
    madoTrayStatusItem.visible = YES;
    madoTrayStatusItem.button.enabled = YES;
    madoTrayStatusItem.button.hidden = NO;
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
  });
}

void MadoTrayHide(void) {
  dispatch_async(dispatch_get_main_queue(), ^{
    if (madoTrayStatusItem == nil) {
      return;
    }

    [[NSStatusBar systemStatusBar] removeStatusItem:madoTrayStatusItem];
    [madoTrayStatusItem release];
    madoTrayStatusItem = nil;
    madoTrayShowItem = nil;
    madoTrayExitItem = nil;
  });
}

void MadoTraySetLocale(char* locale) {
  dispatch_async(dispatch_get_main_queue(), ^{
    if (locale == NULL) {
      madoTrayLocale = @"en";
    } else {
      NSString *value = [NSString stringWithUTF8String:locale];
      madoTrayLocale = [value isEqualToString:@"es"] ? @"es" : @"en";
    }
    MadoTrayUpdateMenuTexts();
  });
}
